import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const laptop = el.querySelector("g");
  animate(laptop, { scaleY: [0.25, 1.1, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hinge = el.querySelector('path[d^="M20.054"]');
  animate(hinge, { scaleX: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
