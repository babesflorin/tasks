<?php

namespace App\Presentation\Transformer;

use App\Domain\Dto\TaskDto;
use League\Fractal\TransformerAbstract;

class TaskTransformer extends TransformerAbstract
{
    public function transform(TaskDto $taskDto)
    {
        return [
            'id' => $taskDto->id,
            'name' => $taskDto->name,
            'description' => $taskDto->description,
            'when' => $taskDto->when,
            'done' => (bool)$taskDto->done,
            'created_at' => $taskDto->createdAt,
            'updated_at' => $taskDto->createdAt,
        ];
    }
}
