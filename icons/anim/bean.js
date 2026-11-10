import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bean = el.querySelector('path[d^="M10.165"]');
  animate(bean, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const seam = el.querySelector('path[d^="M5.341"]');
  animate(seam, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
