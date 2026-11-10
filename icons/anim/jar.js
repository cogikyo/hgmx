import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const jar = el.querySelector('path[d^="M5 3"]');
  animate(jar, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lid = el.querySelector('path[d="M4 3h16"]');
  animate(lid, { translateY: [0, -2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
