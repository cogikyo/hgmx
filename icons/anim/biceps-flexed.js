import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const arm = el.querySelector("g");
  animate(arm, { scale: [1, 1.1, 0.95, 1] }, { duration: τmid });

  const bicep = el.querySelector('path[d^="M15 14"]');
  animate(bicep, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
