import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const lever = el.querySelector('path[d="M18 3h-3"]');
  animate(lever, { rotate: [0, 25, 0] }, { duration: τmid });

  const extinguisher = el.querySelector("g");
  animate(extinguisher, { rotate: [0, -4, 2, 0] }, { duration: τmid, delay: τmin });
}
