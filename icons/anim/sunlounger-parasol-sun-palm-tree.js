import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const canopy = el.querySelector('path[d="m4 14 3-5 5 3Z"]');
  animate(canopy, { scaleX: [0.3, 1.15, 1] }, { duration: τmid });

  const palm = el.querySelector("g");
  animate(palm, { rotate: [0, -10, 6, 0] }, { duration: τmid, delay: τmin });
}
