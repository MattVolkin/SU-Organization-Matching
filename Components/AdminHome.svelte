<script>
    let maxClubsPrePage = $state(10)
</script>
<div class="admin-home">
    <AdminSwitch enabled={true} value={adminPreviewType} onChange={(nextView) => adminPreviewType = nextView} />
    <Header userType="admin" previewAs={adminPreviewType} />    
    <h1>Admin Home</h1>
</div>
