import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const web = el.querySelectorAll('path:not([d^="m16 18"])');
  animate(web, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const data = el.querySelector('path[d^="m16 18"]');
  animate(data, { scale: [0, 1.15, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
