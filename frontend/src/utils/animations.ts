import { SetWindowSize } from "./../../wailsjs/go/main/App";

import type Menu from "./../App.svelte"

let workingarea: HTMLElement;

export function setWorkingArea(element: HTMLElement) {
  workingarea = element;
}

let currentTransitionHandler: (() => void) | null = null;

export function sizeAnimation(
  oldWidth: number, 
  oldHeight: number, 
  newWidth: number, 
  newHeight: number
) {
  if (!workingarea) {
    console.error('Working area not set');
    return;
  }
  
  if (currentTransitionHandler) {
    workingarea.removeEventListener('transitionend', currentTransitionHandler);
    currentTransitionHandler = null;
  }

  let BiggestWidth = Math.max(oldWidth, newWidth);
  let BiggestHeight = Math.max(oldHeight, newHeight);
  
  workingarea.style.transition = "";
  workingarea.style.width = `${oldHeight}px`;
  workingarea.style.height = `${oldHeight}px`;

  SetWindowSize(Math.round(BiggestWidth), Math.round(BiggestHeight));

  const workingAreaTransitionEnd = () => {
    if (!workingarea) {
      console.error('Working area not set');
      return;
    }

    workingarea.removeEventListener('transitionend', workingAreaTransitionEnd);
    currentTransitionHandler = null
    
    workingarea.style.transition = "";
    workingarea.style.width = `${newWidth}px`;
    workingarea.style.height = `${newHeight}px`;

    SetWindowSize(Math.floor(newWidth), Math.floor(newHeight));

    return
  }

  workingarea.addEventListener('transitionend', workingAreaTransitionEnd);
  currentTransitionHandler = workingAreaTransitionEnd
  
  workingarea.style.transition = "width 0.3s, height 0.3s";
  workingarea.style.width = `${newWidth}px`;
  workingarea.style.height = `${newHeight}px`;
}


export function getWorkingArea() {

  return workingarea

}

export function getWorkingAreaSize() {

  if (!workingarea) {
    console.error('Working area not set');
    return;
  }
  
  const computedStyle = workingarea.getBoundingClientRect();
  const currentWidth = Math.round(computedStyle.width);
  const currentHeight = Math.round(computedStyle.height);

  return {
    x: currentWidth,
    y: currentHeight
  }

}