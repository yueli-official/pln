<template>
  <div class="min-h-screen bg-background text-foreground flex flex-col">
    <HeaderArea></HeaderArea>

    <main class="flex-1 overflow-y-scroll">
      <RouterView />
    </main>
    <ToastContainer></ToastContainer>
  </div>
</template>

<script setup lang="ts">
import { ToastContainer } from '@yuelioi/toast'
</script>
