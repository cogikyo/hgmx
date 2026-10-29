import { animate, stagger } from "../motion.js";

export default function (el) {
  const [top, middle, bottom] = el.querySelectorAll("path");
  animate([top, bottom], { translateX: [-4, 1, 0] }, { duration: 0.3, delay: stagger(0.12) });
  animate(middle, { translateX: [4, -1, 0] }, { duration: 0.3, delay: 0.06 });
}
