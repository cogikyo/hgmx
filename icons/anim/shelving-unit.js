import { animate, stagger } from "../motion.js";

export default function (el) {
  const boxes = el.querySelectorAll('path[d^="M12 12V9"], path[d^="M16 20v-3"]');
  animate(boxes, { translateY: [0, -2.5, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
