import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const drops = [el.querySelector('path[d^="M7 16.3"]'), el.querySelector('path[d^="M12.56"]')];
  animate(drops, { translateY: [-4, 1, 0], scaleY: [1.15, 0.85, 1], scaleX: [0.85, 1.15, 1] }, { duration: τmax, delay: stagger(0.06) });
}
