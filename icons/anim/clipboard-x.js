import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector('path[d^="M16 4"]');
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = el.querySelectorAll('path[d^="m14.5"], path[d^="m9.5"]');
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmin, delay: stagger(τmin, { startDelay: τmin }) });

  const clip = el.querySelector("rect");
  animate(clip, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}
