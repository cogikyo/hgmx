import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pillows = el.querySelector("g");
  animate(pillows, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmax });

  const seam = el.querySelector('path[d="M12 4v6"]');
  animate(seam, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
