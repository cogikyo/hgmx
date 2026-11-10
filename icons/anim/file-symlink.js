import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M4 11V4"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const head = el.querySelector('path[d^="m10 18"]');
  animate(head, { translateX: [-4, 2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
