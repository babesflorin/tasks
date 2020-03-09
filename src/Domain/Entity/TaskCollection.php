<?php

namespace App\Domain\Entity;

class TaskCollection extends Collection
{
    protected function getType(): string
    {
        return Task::class;
    }
}