import { animate, stagger } from "../motion.js";

export default function (el) {
  const [top, middle, last] = el.querySelectorAll("path");
  animate(top, { translateX: [-4, 1, 0] }, { duration: 0.3 });
  animate([middle, last], { translateX: [4, -1, 0] }, { duration: 0.3, delay: stagger(0.06, { startDelay: 0.06 }) });
}
