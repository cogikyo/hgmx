import { animate, stagger } from "../motion.js";

export default function (el) {
  const pins = [4, 8, 12, 16, 20].map((x) => el.querySelector(`path[d="M${x} 18v-2"]`));
  animate(pins, { scaleY: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
