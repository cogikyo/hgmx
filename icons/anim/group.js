import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const press = { duration: τmid };
  animate(el.querySelector('rect[x="7"]'), { translateX: [0, 0.6, 0], translateY: [0, 0.6, 0] }, press);
  animate(el.querySelector('rect[x="10"]'), { translateX: [0, -0.6, 0], translateY: [0, -0.6, 0] }, press);

  const corners = el.querySelector("g");
  animate(corners, { scale: [1, 0.92, 1] }, { duration: τmid, delay: τmin });
}
