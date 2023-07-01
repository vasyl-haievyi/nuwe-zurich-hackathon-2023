resource "aws_iam_role" "inserter_lambda_role" {
  name               = "inserter_lambda_iam_role"
  assume_role_policy = templatefile("${path.module}/resources/lambda_role_policy.json", {})
}

resource "aws_iam_policy" "iam_policy_for_inserter_lambda" {
  name = "inserter_lambda_iam_policy"
  policy = templatefile("${path.module}/resources/inserter_policy.json", {
    data_bucket_arn = aws_s3_bucket.data_bucket.arn
    data_table_arn  = aws_dynamodb_table.data_table.arn
  })
}

resource "aws_iam_role_policy_attachment" "attach_iam_policy_to_iam_role_inserter" {
  role       = aws_iam_role.inserter_lambda_role.name
  policy_arn = aws_iam_policy.iam_policy_for_inserter_lambda.arn
}

resource "null_resource" "build_lambda" {
  provisioner "local-exec" {
    command = "cd ${path.module}/../lambda/ && go build  -o ./bin/ ./..."
  }
}

data "archive_file" "zip_the_inserter_code" {
  type        = "zip"
  source_file = "${path.module}/../lambda/bin/inserter"
  output_path = "${path.module}/../build/inserter.zip"

  depends_on = [null_resource.build_lambda]
}

resource "aws_lambda_function" "inserter_function" {
  filename      = "${path.module}/../build/inserter.zip"
  source_code_hash = filebase64sha256("${path.module}/../build/inserter.zip")
  function_name = "inserter_lambda_function"
  role          = aws_iam_role.inserter_lambda_role.arn
  handler       = "inserter"
  runtime       = "go1.x"
  environment {
    variables = {
       DATA_TABLE_NAME = aws_dynamodb_table.data_table.name
    }
  }
  depends_on    = [
    aws_iam_role_policy_attachment.attach_iam_policy_to_iam_role_inserter,
    data.archive_file.zip_the_inserter_code
  ]
}

resource "aws_cloudwatch_log_group" "inserter_function_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.inserter_function.function_name}"
  retention_in_days = 3
}
