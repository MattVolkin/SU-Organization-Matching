<script>
    import Header from "../header.svelte";
    import Footer from "../footer.svelte";
    import DeleteAccount from "../deleteAccount.svelte";
</script>

<div class="delete-page">
    <Header />
    <main class="delete-main">
        <DeleteAccount />
    </main>
    <Footer />
</div>

<style>
    :global(html),
    :global(body) {
        margin: 0;
        padding: 0;
        height: 100%;
    }

    :global(body) {
        overflow: hidden;
    }

    .delete-page {
        height: 100vh;
        display: flex;
        flex-direction: column;
    }

    .delete-main {
        flex: 1;
        min-height: 0;
        display: flex;
    }

    /* .delete-main :global(.swipedContainer) {
        width: 100%;
        height: 100%;
    } */
</style>