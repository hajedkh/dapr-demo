# main.tf

provider "aws" {
  region = "eu-west-1"  # Change to your desired region
}

provider "kubernetes" {
  host                   = aws_eks_cluster.eks_cluster.endpoint
  token                  = data.aws_eks_cluster_auth.cluster_auth.token
  cluster_ca_certificate = base64decode(aws_eks_cluster.eks_cluster.certificate_authority[0].data)
}

data "aws_caller_identity" "current" {}
data "aws_region" "current" {}

# Create a custom VPC
resource "aws_vpc" "custom_vpc" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "custom-vpc"
  }
}

# Create two public subnets in different Availability Zones
resource "aws_subnet" "public_subnet_a" {
  vpc_id                  = aws_vpc.custom_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "eu-west-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "public-subnet-a"
  }
}

resource "aws_subnet" "public_subnet_b" {
  vpc_id                  = aws_vpc.custom_vpc.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "eu-west-1b"
  map_public_ip_on_launch = true

  tags = {
    Name = "public-subnet-b"
  }
}

# Create an Internet Gateway for the VPC
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.custom_vpc.id

  tags = {
    Name = "internet-gateway"
  }
}

# Create a route table
resource "aws_route_table" "public_route_table" {
  vpc_id = aws_vpc.custom_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "public-route-table"
  }
}

# Associate the route table with both public subnets
resource "aws_route_table_association" "public_subnet_a_association" {
  subnet_id      = aws_subnet.public_subnet_a.id
  route_table_id = aws_route_table.public_route_table.id
}

resource "aws_route_table_association" "public_subnet_b_association" {
  subnet_id      = aws_subnet.public_subnet_b.id
  route_table_id = aws_route_table.public_route_table.id
}

# EKS Cluster Resource
resource "aws_eks_cluster" "eks_cluster" {
  name     = "eks-poc-cluster"
  role_arn = aws_iam_role.eks_cluster_role.arn

  vpc_config {
    subnet_ids = [
      aws_subnet.public_subnet_a.id, 
      aws_subnet.public_subnet_b.id  # Use both public subnets
    ]
  }

  depends_on = [aws_iam_role_policy_attachment.cluster_full_permissions]
}

# IAM Role for the EKS Cluster
resource "aws_iam_role" "eks_cluster_role" {
  name = "eks-cluster-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "eks.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Attach the Full Permissions Policy to the EKS Cluster Role
resource "aws_iam_role_policy_attachment" "cluster_full_permissions" {
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
  role       = aws_iam_role.eks_cluster_role.name
}

# IAM Instance Profile for Worker Nodes
resource "aws_iam_instance_profile" "worker_node_profile" {
  name = "worker-node-instance-profile"
  role = aws_iam_role.worker_node_role.name
}

# IAM Role for Worker Nodes
resource "aws_iam_role" "worker_node_role" {
  name = "eks-worker-node-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Attach the Full Permissions Policy to the Worker Node Role
resource "aws_iam_role_policy_attachment" "worker_node_full_permissions" {
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
  role       = aws_iam_role.worker_node_role.name
}

# EKS Node Group
resource "aws_eks_node_group" "eks_node_group" {
  cluster_name    = aws_eks_cluster.eks_cluster.name
  node_group_name = "eks-poc-node-group"
  node_role_arn   = aws_iam_role.worker_node_role.arn
  subnet_ids      = [
    aws_subnet.public_subnet_a.id, 
    aws_subnet.public_subnet_b.id  # Use both public subnets
  ]

  scaling_config {
    desired_size = 2
    max_size     = 3  # You can adjust this
    min_size     = 1
  }

  instance_types = ["t3.small"]

  depends_on = [
    aws_iam_role_policy_attachment.worker_node_full_permissions
  ]
}


# Data block to get the kubeconfig output
data "aws_eks_cluster_auth" "cluster_auth" {
  name = aws_eks_cluster.eks_cluster.name
}


# Output the EKS Cluster Endpoint
output "eks_cluster_endpoint" {
  value = aws_eks_cluster.eks_cluster.endpoint
}

# Output the command to configure kubectl
output "connect_command" {
  value = "aws eks --region ${data.aws_region.current.name} update-kubeconfig --name ${aws_eks_cluster.eks_cluster.name}"
}
