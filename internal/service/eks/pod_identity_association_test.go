// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eks_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	tfeks "github.com/hashicorp/terraform-provider-aws/internal/service/eks"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccEKSPodIdentityAssociation_basic(t *testing.T) {
	ctx := acctest.Context(t)

	var podidentityassociation eks.DescribePodIdentityAssociationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_pod_identity_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.EKSEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EKSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPodIdentityAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPodIdentityAssociationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPodIdentityAssociationExists(ctx, resourceName, &podidentityassociation),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace"),
					resource.TestCheckResourceAttrSet(resourceName, "role_arn"),
					resource.TestCheckResourceAttrSet(resourceName, "service_account"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccEKSPodIdentityAssociation_disappears(t *testing.T) {
	ctx := acctest.Context(t)

	var podidentityassociation eks.DescribePodIdentityAssociationOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_eks_pod_identity_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.EKSEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.EKSEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPodIdentityAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPodIdentityAssociationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPodIdentityAssociationExists(ctx, resourceName, &podidentityassociation),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfeks.ResourcePodIdentityAssociation, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckPodIdentityAssociationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EKSClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_eks_pod_identity_association" {
				continue
			}

			_, err := conn.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
				AssociationId: aws.String(rs.Primary.ID),
				ClusterName:   aws.String(rs.Primary.Attributes["cluster_name"]),
			})
			if errs.IsA[*types.ResourceNotFoundException](err) {
				return nil
			}
			if err != nil {
				return create.Error(names.EKS, create.ErrActionCheckingDestroyed, tfeks.ResNamePodIdentityAssociation, rs.Primary.ID, err)
			}

			return create.Error(names.EKS, create.ErrActionCheckingDestroyed, tfeks.ResNamePodIdentityAssociation, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckPodIdentityAssociationExists(ctx context.Context, name string, podidentityassociation *eks.DescribePodIdentityAssociationOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.EKS, create.ErrActionCheckingExistence, tfeks.ResNamePodIdentityAssociation, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.EKS, create.ErrActionCheckingExistence, tfeks.ResNamePodIdentityAssociation, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EKSClient(ctx)
		resp, err := conn.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
			AssociationId: aws.String(rs.Primary.ID),
			ClusterName:   aws.String(rs.Primary.Attributes["cluster_name"]),
		})

		if err != nil {
			return create.Error(names.EKS, create.ErrActionCheckingExistence, tfeks.ResNamePodIdentityAssociation, rs.Primary.ID, err)
		}

		*podidentityassociation = *resp

		return nil
	}
}

func testAccCheckPodIdentityAssociationNotRecreated(before, after *eks.DescribePodIdentityAssociationOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.ToString(before.Association.AssociationId), aws.ToString(after.Association.AssociationId); before != after {
			return create.Error(names.EKS, create.ErrActionCheckingNotRecreated, tfeks.ResNamePodIdentityAssociation, before, errors.New("recreated"))
		}

		return nil
	}
}

func testAccPodIdentityAssociationConfig_clusterBase(rName string) string {
	return acctest.ConfigCompose(
		acctest.ConfigAvailableAZsNoOptInDefaultExclude(),
		fmt.Sprintf(`
resource "aws_iam_role" "test" {
  name_prefix = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Action": [
        "sts:AssumeRole",
        "sts:TagSession"
      ]
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "test_cluster" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.test.name
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name                          = %[1]q
    "kubernetes.io/cluster/%[1]s" = "shared"
  }
}

resource "aws_subnet" "test" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = cidrsubnet(aws_vpc.test.cidr_block, 8, count.index)
  vpc_id            = aws_vpc.test.id

  tags = {
    Name                          = %[1]q
    "kubernetes.io/cluster/%[1]s" = "shared"
  }
}

resource "aws_eks_cluster" "test" {
  name     = %[1]q
  role_arn = aws_iam_role.test.arn

  vpc_config {
    subnet_ids = aws_subnet.test[*].id
  }

  depends_on = [aws_iam_role_policy_attachment.test_cluster]
}
`, rName))
}

func testAccPodIdentityAssociationConfig_podIdentityRoleBase(rName string) string {
	return acctest.ConfigCompose(
		fmt.Sprintf(`
resource "aws_iam_role" "example" {
  name = %[1]q

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "pods.eks.amazonaws.com"
      },
      "Action": [
        "sts:AssumeRole",
        "sts:TagSession"
      ]
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "example" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
  role       = aws_iam_role.example.name
}
`, rName))
}
func testAccPodIdentityAssociationConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccPodIdentityAssociationConfig_clusterBase(rName),
		testAccPodIdentityAssociationConfig_podIdentityRoleBase(rName),
		fmt.Sprintf(`
resource "aws_eks_pod_identity_association" "test" {
  cluster_name    = aws_eks_cluster.test.name
  namespace       = %[1]q
  service_account = "%[1]s-sa"
  role_arn        = aws_iam_role.test.arn
}
`, rName))
}
