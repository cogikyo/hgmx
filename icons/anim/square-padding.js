import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const top = el.querySelector('path[d="M14 7h-4"]');
  animate(top, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bottom = el.querySelector('path[d="M14 17h-4"]');
  animate(bottom, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const left = el.querySelector('path[d="M7 14v-4"]');
  animate(left, { translateX: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const right = el.querySelector('path[d="M17 14v-4"]');
  animate(right, { translateX: [4, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
