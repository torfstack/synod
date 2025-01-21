<script lang="ts">
    import api from "$lib/api";
    import cfg from "$lib/config";
    import {Button, Checkbox, Helper, Img, Input, Label} from "flowbite-svelte";
    import {EnvelopeSolid} from "flowbite-svelte-icons";
    import KayHeader from "../components/KayHeader.svelte";

    let email = $state("");
    let password = $state("");

    function signInWithProvider() {
        window.open(cfg.backendAuthStartUrl, "_self");
    }
</script>

<div>
    <KayHeader text="Login or register, keep "/>
    <div class="p-3">

        <br>

        <div class="flex justify-center">
            <div class="lg:w-1/2 2xl:w-1/4 p-6 rounded bg-gradient-to-r dark:from-sky-800 from-sky-500 dark:to-gray-900 to-gray-600">
                <form>
                    <div class="mb-3">
                        <Button color="light" on:click={signInWithProvider}>
                            <Img src="https://img.icons8.com/color/16/000000/google-logo.png"/>Login with Google
                        </Button>
                    </div>
                    <div class="mb-3">
                        <Label class="space-y-2">
                            <span>Email address</span> <Input bind:value={email} placeholder="kay@vault.com" size="md" type="email">
                                {#snippet left()}
                                    <EnvelopeSolid name="envelope-solid" />
                                {/snippet}
                            </Input>
                        </Label>
                    </div>
                    <div class="mb-3">
                        <Label class="space-y-2">
                            <span>Password</span>
                            <Input bind:value={password} placeholder="much secret" size="md" type="password"></Input>
                        </Label>
                    </div>
                    <div class="mb-3">
                        <Checkbox>Remember me</Checkbox>
                    </div>
                    <div class="grid lg:grid-cols-2 gap-4">
                        <Button>Login</Button>
                        <Button color="light">Register</Button>
                    </div>
                </form>
            </div>
        </div>

    </div>
</div>
