import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const boat = el.querySelector("g");
  animate(boat, { rotate: [0, -8, 6, -3, 0] }, { duration: τmax });

  const sail = el.querySelector('path[d^="M9.159"]');
  animate(sail, { scaleX: [1, 1.12, 1] }, { duration: τmid, delay: τmin });
}
