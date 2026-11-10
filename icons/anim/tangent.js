import { animate, τmid } from "../motion.js";

export default function (el) {
  const curve = el.querySelector('path[d^="M12 22"]');
  animate(curve, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
}
