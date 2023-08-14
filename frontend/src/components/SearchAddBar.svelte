<script lang="ts">
    import {ButtonGroup, Input, InputAddon} from "flowbite-svelte";
    import {Icon} from "flowbite-svelte-icons";
    export let filterValue: string;
    export let clickedAdd: boolean;

    function didClickAdd() {
        clickedAdd = false
        clickedAdd = true
    }
</script>

<div class="flex justify-center">
    <div class="lg:w-1/2">
        <ButtonGroup class="flex">
            <InputAddon>
                <button on:click={didClickAdd}>
                    <Icon name="plus-outline" data-bs-toggle="modal" data-bs-target="#exampleModal" />
                </button>
            </InputAddon>
            <Input bind:value={filterValue} type="text" size="lg" placeholder="Add/Search Secrets" name="New Secret">
            </Input>
        </ButtonGroup>
    </div>
</div>