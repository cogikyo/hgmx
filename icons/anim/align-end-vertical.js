import { animate, τmax } from "../motion.js";

export default function (el) {
  const long = el.querySelector('rect[width="16"]');
  animate(long, { translateX: [-1.5, 0.5, 0] }, { duration: τmax });

  const short = el.querySelector('rect[width="9"]');
  animate(short, { translateX: [-3.5, 0.5, 0] }, { duration: τmax });
}
