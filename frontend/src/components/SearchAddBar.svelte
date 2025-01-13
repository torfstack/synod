<script lang="ts">
    import {ButtonGroup, Input, InputAddon} from "flowbite-svelte";
    import { PlusOutline } from "flowbite-svelte-icons";
    interface Props {
        filterValue: string;
        clickedAdd: boolean;
    }

    let { filterValue = $bindable(), clickedAdd = $bindable() }: Props = $props();

    function didClickAdd() {
        clickedAdd = false
        clickedAdd = true
    }
</script>

<div class="flex justify-center">
    <div class="lg:w-1/2">
        <ButtonGroup class="flex">
            <InputAddon>
                <button onclick={didClickAdd}>
                    <PlusOutline data-bs-target="#exampleModal" data-bs-toggle="modal" name="plus-outline" />
                </button>
            </InputAddon>
            <Input bind:value={filterValue} name="New Secret" placeholder="Add/Search Secrets" size="lg" type="text">
            </Input>
        </ButtonGroup>
    </div>
</div>