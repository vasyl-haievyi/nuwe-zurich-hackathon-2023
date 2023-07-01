resource "aws_s3_bucket" "data_bucket" {
  bucket        = "data-bucket"
  force_destroy = true
}

resource "aws_s3_bucket_notification" "data_bucket_to_inserter" {
  bucket = aws_s3_bucket.data_bucket.id
  lambda_function {
    lambda_function_arn = aws_lambda_function.inserter_function.arn
    events              = ["s3:ObjectCreated:*"]
  }

  depends_on = [ aws_lambda_permission.data_bucket_invoke_inserter ]
}

resource "aws_lambda_permission" "data_bucket_invoke_inserter" {
  statement_id  = "AllowS3Invoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.inserter_function.function_name
  principal     = "s3.amazonaws.com"
  source_arn    = aws_s3_bucket.data_bucket.arn
}