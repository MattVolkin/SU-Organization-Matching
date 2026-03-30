<script>
import AdminSwitch from './AdminSwitch.svelte'
import Header from './header.svelte'
import Footer from './footer.svelte'
import { APICreater } from './APIHandler.svelte'
let clubName = $state('');
let APIobj = $state({});
APIobj.APICreater('POST', '/api/createClub', {"name": `${clubName}`}) ;

</script>
<API bind:this={APIobj} />

