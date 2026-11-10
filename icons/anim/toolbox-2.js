import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = el.querySelector('path[d^="M4 21"]');
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const handle = el.querySelector('path[d^="M8 7"]');
  animate(handle, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
