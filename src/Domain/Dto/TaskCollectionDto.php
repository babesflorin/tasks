<?php


namespace App\Domain\Dto;


use App\Domain\Entity\Collection;

class TaskCollectionDto extends Collection
{
    protected function getType(): string
    {
        return TaskDto::class;
    }
}