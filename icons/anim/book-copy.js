import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const front = el.querySelector('path[d^="M9 15"]');
  animate(front, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const back = el.querySelectorAll('path[d^="M5 7"], path[d^="M5.803"]');
  animate(back, { translateX: [4, -1, 0], translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
