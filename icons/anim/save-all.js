import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const disk = el.querySelectorAll('path[d^="M8 18"], path[d^="M18 18"]');
  animate(disk, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const behind = el.querySelector('path[d^="M18 22"]');
  animate(behind, { translateX: [4, -1, 0], translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
