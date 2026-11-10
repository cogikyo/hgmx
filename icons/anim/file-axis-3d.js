import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const axes = el.querySelector('path[d="M8 10v8h8"]');
  animate(axes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const depth = el.querySelector('path[d="m8 18 4-4"]');
  animate(depth, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
