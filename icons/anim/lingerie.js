import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bra = el.querySelector("g");
  animate(bra, { scale: [1, 1.2, 0.95, 1] }, { duration: τmid });

  const brief = el.querySelector('path[d^="M3 15"]');
  animate(brief, { scale: [1, 1.2, 0.95, 1] }, { duration: τmid, delay: τmin });
}
