<?php

namespace App\Domain\Transformer;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Entity\TaskCollection;

class TaskCollectionDtoToCollectionTransformer
{
    /**
     * @var TaskDtoToEntityTransformer
     */
    private $entityTransformer;

    public function __construct(TaskDtoToEntityTransformer $entityTransformer)
    {
        $this->entityTransformer = $entityTransformer;
    }

    public function transform(TaskCollectionDto $taskCollectionDto): TaskCollection
    {
        $taskCollection = new TaskCollection();
        foreach ($taskCollectionDto as $taskDto) {
            $taskCollection->append($this->entityTransformer->transform($taskDto));
        }

        return $taskCollection;
    }

    public function reverseTransform(TaskCollection $taskCollection)
    {
        $taskCollectionDto = new TaskCollectionDto();
        foreach ($taskCollection as $task) {
            $taskCollectionDto->append($this->entityTransformer->reverseTransform($task));
        }

        return $taskCollectionDto;
    }
}
