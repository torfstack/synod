<script lang="ts">
    import {config} from "$lib/config";
    import {Button, Img} from "flowbite-svelte";
    import KayHeader from "../components/KayHeader.svelte";

    let email = $state("");
    let password = $state("");

    function signInWithProvider() {
        window.open(config.backendAuthStartUrl, "_self");
    }
</script>

<div>
    <KayHeader text="Login, keep "/>
    <div class="p-3">
        <br>
        <div class="flex justify-center">
            <form>
                <div class="mb-3 flex justify-center">
                    <Button color="light" on:click={signInWithProvider}>
                        <Img src="https://img.icons8.com/color/16/000000/google-logo.png"/>&nbsp;Login with Google
                    </Button>
                </div>
            </form>
        </div>
    </div>
</div>
