import { animate, stagger } from "../motion.js";

export default function (el) {
  const lines = [7, 12, 17].map((y) => el.querySelector(`path[d="M10 ${y}h4"]`));
  animate(lines, { scaleX: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
