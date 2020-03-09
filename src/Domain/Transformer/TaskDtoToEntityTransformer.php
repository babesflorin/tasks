<?php

namespace App\Domain\Transformer;

use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;

class TaskDtoToEntityTransformer
{
    public function transform(TaskDto $taskDto)
    {
        return new Task($taskDto->name, $taskDto->description, \DateTime::createFromFormat('Y-m-d', $taskDto->when)->setTime(0,0,0));
    }

    public function reverseTransform(Task $task)
    {
        $taskDto = new TaskDto();
        $taskDto->name = $task->getName();
        $taskDto->description = $task->getDescription();
        $taskDto->when = $task->getWhen()->format('Y-m-d');
        $taskDto->done = $task->isDone();
        $taskDto->createdAt = $task->getCreatedAt();
        $taskDto->updatedAt = $task->getUpdateAt();

        return $taskDto;
    }
}