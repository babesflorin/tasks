<?php

namespace App\Domain\Validator;

use App\Domain\Dto\TaskDto;
use App\Domain\Exception\InvalidTaskException;

class TaskValidator
{
    /**
     * @throws InvalidTaskException
     */
    public function validate(TaskDto $taskDto, bool $shouldHaveId = false)
    {
        $errors = [];
        if (!is_string($taskDto->name) || empty($taskDto->name)) {
            $errors[] = 'Task name is not valid!';
        }
        if ($shouldHaveId) {
            if (null === $taskDto->id) {
                $errors[] = 'We need an id to know which entity to update!';
            } elseif (!is_integer($taskDto->id)) {
                $errors[] = 'The id must be an integer!';
            }
        }

        if (!is_string($taskDto->description) || empty($taskDto->description)) {
            $errors[] = 'Task description is not valid!';
        }

        if (empty($taskDto->when)) {
            $errors[] = 'Task must have a date!';
        } elseif (false === ($when = \DateTime::createFromFormat("Y-m-d", $taskDto->when))) {
            $errors[] = '`when` is not a valid date!';
        } elseif ($when < (new \DateTime())) {
            $errors[] = 'You can\'t do a task in a the past. Or you can?';
        }
        if (!empty($errors)) {
            throw new InvalidTaskException($errors);
        }
    }
}