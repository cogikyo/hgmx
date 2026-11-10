import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const loop = el.querySelectorAll('path[d^="M3.34"], path[d^="M7 16.93"]');
  animate(loop, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const pointer = el.querySelector('path[d^="M14.33"]');
  animate(pointer, { translateX: [0, -2, 1, 0], translateY: [0, -2, 1, 0] }, { duration: τmid, delay: τmin });
}
