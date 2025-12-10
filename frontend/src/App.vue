<template>
  <div class="min-h-screen bg-background text-foreground">
    <HeaderArea></HeaderArea>
    <HomeArea></HomeArea>

    <main>
      <RouterView />
    </main>
    <ToastContainer></ToastContainer>
  </div>
</template>

<script setup lang="ts">
import { ToastContainer } from '@yuelioi/toast'
</script>
