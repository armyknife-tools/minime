// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package templates

// generateAWSTemplates generates templates for AWS resources
func generateAWSTemplates() []Template {
	var templates []Template

	// AWS S3 Bucket
	templates = append(templates, Template{
		Provider:    "aws",
		Resource:    "s3",
		DisplayName: "S3 Bucket",
		Description: "Amazon S3 bucket with all common configuration options",
		Category:    "Storage",
		Tags:        "storage,s3,bucket",
		Content: `# AWS S3 Bucket Template
# This template includes all common properties for an S3 bucket

resource "aws_s3_bucket" "example" {
  bucket = "my-example-bucket"
  
  # Access control
  acl    = "private"  # Options: private, public-read, public-read-write, etc.
  
  # Website configuration
  website {
    index_document = "index.html"
    error_document = "error.html"
    routing_rules = <<EOF
[{
    "Condition": {
        "KeyPrefixEquals": "docs/"
    },
    "Redirect": {
        "ReplaceKeyPrefixWith": "documents/"
    }
}]
EOF
  }
  
  # Versioning configuration
  versioning {
    enabled = true
    mfa_delete = false
  }
  
  # Lifecycle rules
  lifecycle_rule {
    id      = "log-retention"
    enabled = true
    
    prefix  = "logs/"
    
    transition {
      days          = 30
      storage_class = "STANDARD_IA"
    }
    
    transition {
      days          = 90
      storage_class = "GLACIER"
    }
    
    expiration {
      days = 365
    }
  }
  
  # Server-side encryption
  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
  
  # Object lock configuration
  object_lock_configuration {
    object_lock_enabled = "Enabled"
    rule {
      default_retention {
        mode = "GOVERNANCE"
        days = 30
      }
    }
  }
  
  # Logging
  logging {
    target_bucket = "example-logs-bucket"
    target_prefix = "log/"
  }
  
  # CORS configuration
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET", "PUT", "POST"]
    allowed_origins = ["https://example.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
  
  # Tags
  tags = {
    Name        = "My Example Bucket"
    Environment = "Dev"
    Project     = "Example"
  }
}

# S3 bucket policy
resource "aws_s3_bucket_policy" "example" {
  bucket = aws_s3_bucket.example.id
  
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "${aws_s3_bucket.example.arn}/*"
    }
  ]
}
POLICY
}

# S3 bucket notification configuration
resource "aws_s3_bucket_notification" "example" {
  bucket = aws_s3_bucket.example.id
  
  topic {
    topic_arn     = "arn:aws:sns:us-east-1:123456789012:example-topic"
    events        = ["s3:ObjectCreated:*"]
    filter_prefix = "logs/"
    filter_suffix = ".log"
  }
  
  queue {
    queue_arn     = "arn:aws:sqs:us-east-1:123456789012:example-queue"
    events        = ["s3:ObjectCreated:*"]
    filter_prefix = "images/"
    filter_suffix = ".jpg"
  }
  
  lambda_function {
    lambda_function_arn = "arn:aws:lambda:us-east-1:123456789012:function:example-function"
    events              = ["s3:ObjectCreated:*", "s3:ObjectRemoved:*"]
    filter_prefix       = "data/"
    filter_suffix       = ".csv"
  }
}
`,
	})

	// AWS EC2 Instance
	templates = append(templates, Template{
		Provider:    "aws",
		Resource:    "ec2",
		DisplayName: "EC2 Instance",
		Description: "Amazon EC2 instance with all common configuration options",
		Category:    "Compute",
		Tags:        "compute,ec2,instance",
		Content: `# AWS EC2 Instance Template
# This template includes all common properties for an EC2 instance

resource "aws_instance" "example" {
  ami           = "ami-0c55b159cbfafe1f0"  # Amazon Linux 2 AMI (HVM), SSD Volume Type
  instance_type = "t2.micro"
  
  # Availability Zone
  availability_zone = "us-west-2a"
  
  # VPC and Subnet
  subnet_id = aws_subnet.example.id
  
  # Security Groups
  vpc_security_group_ids = [aws_security_group.example.id]
  
  # Key Pair for SSH access
  key_name = "example-key-pair"
  
  # EBS Root Volume
  root_block_device {
    volume_type           = "gp2"
    volume_size           = 50
    delete_on_termination = true
    encrypted             = true
    
    tags = {
      Name = "Root Volume"
    }
  }
  
  # Additional EBS Volumes
  ebs_block_device {
    device_name           = "/dev/sdf"
    volume_type           = "gp2"
    volume_size           = 100
    delete_on_termination = true
    encrypted             = true
    
    tags = {
      Name = "Data Volume"
    }
  }
  
  # User Data for instance initialization
  user_data = <<-EOF
              #!/bin/bash
              yum update -y
              yum install -y httpd
              systemctl start httpd
              systemctl enable httpd
              echo "<h1>Hello from OpenTofu</h1>" > /var/www/html/index.html
              EOF
  
  # IAM Instance Profile
  iam_instance_profile = aws_iam_instance_profile.example.name
  
  # Detailed monitoring
  monitoring = true
  
  # Disable source/destination checks for NAT instances
  source_dest_check = false
  
  # Tenancy (default, dedicated, or host)
  tenancy = "default"
  
  # Credit specification for T2/T3 instances
  credit_specification {
    cpu_credits = "unlimited"
  }
  
  # Capacity reservation
  capacity_reservation_specification {
    capacity_reservation_preference = "open"
  }
  
  # Placement group
  placement_group = aws_placement_group.example.id
  
  # Tags
  tags = {
    Name        = "Example EC2 Instance"
    Environment = "Development"
    Project     = "OpenTofu Demo"
  }
}

# Security Group for EC2 Instance
resource "aws_security_group" "example" {
  name        = "example-security-group"
  description = "Allow inbound traffic for the example EC2 instance"
  vpc_id      = aws_vpc.example.id
  
  # SSH access
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH access"
  }
  
  # HTTP access
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP access"
  }
  
  # HTTPS access
  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTPS access"
  }
  
  # Outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }
  
  # Tags
  tags = {
    Name = "Example Security Group"
  }
}
`,
	})

	// AWS RDS Database
	templates = append(templates, Template{
		Provider:    "aws",
		Resource:    "rds",
		DisplayName: "RDS Database",
		Description: "Amazon RDS database with all common configuration options",
		Category:    "Database",
		Tags:        "database,rds,mysql,postgres",
		Content: `# AWS RDS Database Template
# This template includes all common properties for an RDS database

resource "aws_db_instance" "example" {
  # Engine options
  engine         = "mysql"  # Options: mysql, postgres, oracle-ee, sqlserver-ee, etc.
  engine_version = "8.0"
  
  # Settings
  identifier        = "example-db"
  instance_class    = "db.t3.medium"
  allocated_storage = 20
  storage_type      = "gp2"
  
  # Credentials
  username = "admin"
  password = "YourStrongPasswordHere"  # In production, use aws_secretsmanager_secret or similar
  
  # Database name
  name = "exampledb"
  
  # Network & Security
  vpc_security_group_ids = [aws_security_group.db.id]
  db_subnet_group_name   = aws_db_subnet_group.example.name
  publicly_accessible    = false
  
  # Database port
  port = 3306
  
  # Database options
  parameter_group_name = aws_db_parameter_group.example.name
  option_group_name    = aws_db_option_group.example.name
  
  # Backup
  backup_retention_period = 7
  backup_window           = "03:00-04:00"
  copy_tags_to_snapshot   = true
  delete_automated_backups = true
  
  # Maintenance
  maintenance_window = "Mon:00:00-Mon:03:00"
  auto_minor_version_upgrade = true
  
  # Deletion protection
  deletion_protection = true
  
  # Enhanced monitoring
  monitoring_interval = 60
  monitoring_role_arn = aws_iam_role.rds_monitoring.arn
  
  # Performance Insights
  performance_insights_enabled = true
  performance_insights_retention_period = 7  # days
  
  # Storage autoscaling
  max_allocated_storage = 100
  
  # Multi-AZ deployment
  multi_az = true
  
  # Storage encryption
  storage_encrypted = true
  kms_key_id        = aws_kms_key.example.arn
  
  # Tags
  tags = {
    Name        = "Example RDS Instance"
    Environment = "Production"
    Project     = "OpenTofu Demo"
  }
}

# DB Subnet Group
resource "aws_db_subnet_group" "example" {
  name       = "example-db-subnet-group"
  subnet_ids = [aws_subnet.private1.id, aws_subnet.private2.id]
  
  tags = {
    Name = "Example DB Subnet Group"
  }
}

# DB Parameter Group
resource "aws_db_parameter_group" "example" {
  name   = "example-db-parameter-group"
  family = "mysql8.0"
  
  parameter {
    name  = "character_set_server"
    value = "utf8"
  }
  
  parameter {
    name  = "character_set_client"
    value = "utf8"
  }
  
  tags = {
    Name = "Example DB Parameter Group"
  }
}

# DB Option Group
resource "aws_db_option_group" "example" {
  name                 = "example-db-option-group"
  engine_name          = "mysql"
  major_engine_version = "8.0"
  
  option {
    option_name = "MARIADB_AUDIT_PLUGIN"
    
    option_settings {
      name  = "SERVER_AUDIT_EVENTS"
      value = "CONNECT,QUERY"
    }
  }
  
  tags = {
    Name = "Example DB Option Group"
  }
}

# Security Group for RDS
resource "aws_security_group" "db" {
  name        = "example-db-security-group"
  description = "Allow inbound traffic to the RDS instance"
  vpc_id      = aws_vpc.example.id
  
  # MySQL/Aurora access
  ingress {
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.example.id]  # Allow access from EC2 security group
    description     = "MySQL access from EC2 instances"
  }
  
  # Outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    description = "Allow all outbound traffic"
  }
  
  tags = {
    Name = "Example DB Security Group"
  }
}
`,
	})

	return templates
}
