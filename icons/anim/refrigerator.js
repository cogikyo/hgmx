import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const fridge = el.querySelector('path[d^="M5 6a4"]');
  animate(fridge, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const handle = el.querySelector('path[d="M15 7v6"]');
  animate(handle, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
