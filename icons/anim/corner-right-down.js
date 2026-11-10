import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shaft = el.querySelector('path[d^="M4 4"]');
  animate(shaft, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const head = el.querySelector('path[d^="m10 15"]');
  animate(head, { translateY: [-4, 2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
