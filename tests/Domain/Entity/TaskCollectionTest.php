<?php

namespace App\Tests\Domain\Entity;

use App\Domain\Dto\TaskCollectionDto;
use App\Domain\Dto\TaskDto;
use App\Domain\Entity\Task;
use App\Domain\Entity\TaskCollection;
use App\Tests\TestCase;

class TaskCollectionTest extends TestCase
{
    public function testInvalidAppend()
    {
        $taskCollection = new TaskCollection();
        $this->expectException(\InvalidArgumentException::class);
        $taskCollection->append(new \stdClass());
    }

    public function testConstructNotIterable()
    {
        $this->expectException(\InvalidArgumentException::class);
        $taskCollection = new TaskCollection(new \stdClass());
    }
}
