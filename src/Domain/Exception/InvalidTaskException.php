<?php

namespace App\Domain\Exception;

class InvalidTaskException extends ValidationException
{
    public function __construct($errors)
    {
        parent::__construct("Task is not valid!", $errors);
    }
}