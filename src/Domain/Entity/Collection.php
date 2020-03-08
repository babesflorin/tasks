<?php

namespace App\Domain\Entity;

abstract class Collection extends \ArrayObject
{
    public function __construct($input = [], $flags = 0, $iterator_class = "ArrayIterator")
    {
        if (is_iterable($input)) {
            foreach ($input as $value) {
                $this->checkType($value);
            }
        } else {
            $this->checkType($input);
        }
        parent::__construct(
            $input,
            $flags,
            $iterator_class
        );
    }

    public function append($value)
    {
        $this->checkType($value);
        parent::append($value);
    }

    public function offsetSet($index, $newval)
    {
        $this->checkType($newval);
        parent::offsetSet($index, $newval);
    }

    private function checkType($value)
    {
        $collectionType = $this->getType();
        if (!$value instanceof $collectionType) {
            throw new \InvalidArgumentException("You must provide a ".Task::class." object!");
        }
    }

    abstract protected function getType(): string;
}