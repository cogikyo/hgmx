import { animate, τmid } from "../motion.js";

export default function (el) {
  const beacon = el.querySelector('path[d="M13 2h-2"]');
  animate(beacon, { scale: [1, 1.75, 1] }, { duration: τmid });
}
