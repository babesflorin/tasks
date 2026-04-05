<?php

namespace App\Domain\Exception;

class ValidationException extends \Exception
{
    protected $errors;

    public function __construct($message, $errors)
    {
        parent::__construct(
            $message
        );
        $this->errors = $errors;
    }

    public function getErrors(): array
    {
        return $this->errors;
    }
}
